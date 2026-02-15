package db

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	// Security limits to prevent DoS attacks.
	MaxFilters       = 10       // Maximum number of filters allowed
	MaxSortFields    = 5        // Maximum number of sort fields allowed
	MaxArrayElements = 50       // Maximum elements in IN/NIN arrays
	MaxPageNumber    = 1000000  // Maximum page number to prevent overflow
	MaxOffset        = 10000000 // Maximum offset to prevent excessive queries
	MaxInputLength   = 100      // Maximum length for string inputs
	MaxLimit         = 1000     // Maximum limit per page
	MinLimit         = 1        // Minimum limit per page

	DefaultLimit = 500 // Default limit per page
	DefaultPage  = 1   // Default page number
)

var (
	ErrInvalidOperator = errors.New("invalid operator")
	ErrInvalidField    = errors.New("invalid field")
	ErrInvalidValue    = errors.New("invalid value")
)

// FieldMeta stores validator and optional DB column.
type FieldMeta struct {
	Validator FieldValidator
	Column    string
}

// QueryFilter represents a single filter condition.
type QueryFilter struct {
	Field    string
	Operator string
	Value    any
}

// QueryOptions contains query parameters.
type QueryOptions struct {
	Filters []QueryFilter
	Sort    []string
	Page    int
	Limit   int
}

// OperatorHandler applies a filter to a GORM query.
type OperatorHandler func(db *gorm.DB, field string, value any) *gorm.DB

type FieldValidator func(value string) (any, error)

// QueryBuilder builds dynamic GORM queries.
type QueryBuilder struct {
	fields    map[string]FieldMeta
	operators map[string]OperatorHandler
}

// NewQueryBuilder creates a new QueryBuilder.
func NewQueryBuilder() *QueryBuilder {
	qb := &QueryBuilder{
		fields:    make(map[string]FieldMeta),
		operators: make(map[string]OperatorHandler),
	}
	qb.registerDefaultOperators()
	return qb
}

// RegisterField registers a field with validator and optional DB column.
func (qb *QueryBuilder) RegisterField(field string, validator FieldValidator, column ...string) {
	col := field
	if len(column) > 0 {
		col = column[0]
	}
	qb.fields[field] = FieldMeta{
		Validator: validator,
		Column:    col,
	}
}

// ParseQueryParams parses URL params into QueryOptions.
func (qb *QueryBuilder) ParseQueryParams(query url.Values) (*QueryOptions, error) {
	opts := &QueryOptions{Page: DefaultPage, Limit: DefaultLimit}

	for key, values := range query {
		if len(values) == 0 {
			continue
		}
		value := values[0]

		if err := qb.handleQueryParam(key, value, opts); err != nil {
			return nil, err
		}
	}

	// Prevent integer overflow in offset calculation
	if opts.Limit > 0 {
		offset := (opts.Page - 1) * opts.Limit
		if offset < 0 || offset > MaxOffset {
			return nil, ErrInvalidValue
		}
	}

	return opts, nil
}

func (qb *QueryBuilder) handleQueryParam(key, value string, opts *QueryOptions) error {
	switch strings.ToLower(key) {
	case "page":
		p, err := strconv.Atoi(value)
		if err != nil || p <= 0 || p > MaxPageNumber {
			return ErrInvalidValue
		}
		opts.Page = p
	case "limit":
		l, err := strconv.Atoi(value)
		if err != nil || l < MinLimit || l > MaxLimit {
			return ErrInvalidValue
		}
		opts.Limit = l
	case "sort":
		if len(opts.Sort) >= MaxSortFields {
			return ErrInvalidValue
		}
		fields := strings.Split(value, ",")
		for _, f := range fields {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			if len(opts.Sort) >= MaxSortFields {
				return ErrInvalidValue
			}
			col := f
			if strings.HasPrefix(f, "-") {
				col = f[1:]
			}
			if _, ok := qb.fields[col]; !ok {
				return ErrInvalidField
			}
			opts.Sort = append(opts.Sort, f)
		}
	default:
		if len(opts.Filters) >= MaxFilters {
			return ErrInvalidValue
		}
		filter, err := qb.parseFilter(key, value)
		if err != nil {
			return err
		}
		opts.Filters = append(opts.Filters, *filter)
	}
	return nil
}

// parseFilter parses a single filter.
func (qb *QueryBuilder) parseFilter(key, value string) (*QueryFilter, error) {
	field := key
	operator := "eq"

	if i := strings.Index(key, "["); i != -1 {
		j := strings.Index(key, "]")
		if j <= i || j != len(key)-1 {
			return nil, ErrInvalidValue
		}
		field = key[:i]
		operator = strings.ToLower(key[i+1 : j])
		if operator == "" {
			return nil, ErrInvalidValue
		}
	}

	_, ok := qb.fields[field]
	if !ok {
		return nil, ErrInvalidField
	}
	if _, ok := qb.operators[operator]; !ok {
		return nil, ErrInvalidOperator
	}

	val, err := qb.convertValue(field, operator, value)
	if err != nil {
		return nil, err
	}

	return &QueryFilter{
		Field:    field,
		Operator: operator,
		Value:    val,
	}, nil
}

// convertValue converts string value to proper type.
func (qb *QueryBuilder) convertValue(field, operator, value string) (any, error) {
	parse := func(v string) (any, error) {
		if fn, ok := qb.fields[field]; ok {
			return fn.Validator(v)
		}
		return v, nil
	}

	switch operator {
	case "in", "nin":
		parts := strings.Split(value, ",")
		out := make([]any, 0, len(parts))
		for _, p := range parts {
			v, err := parse(strings.TrimSpace(p))
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	case "between":
		parts := strings.SplitN(value, ",", 2) //nolint:mnd // between operator requires exactly 2 parts
		if len(parts) != 2 {                   //nolint:mnd // between operator requires exactly 2 parts
			return nil, ErrInvalidValue
		}
		start, err := parse(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}
		end, err := parse(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}
		return []any{start, end}, nil
	case "null", "notnull":
		return nil, nil //nolint:nilnil // intentional for null operators
	default:
		return parse(value)
	}
}

// Scope returns a GORM scope applying filters, sorting, pagination.
func (qb *QueryBuilder) Scope(opts *QueryOptions) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, f := range opts.Filters {
			col := f.Field
			if meta, ok := qb.fields[f.Field]; ok {
				col = meta.Column
			}
			db = qb.operators[f.Operator](db, col, f.Value)
		}

		if len(opts.Sort) > 0 {
			var clauses []string
			for _, s := range opts.Sort {
				desc := ""
				col := s
				if strings.HasPrefix(s, "-") {
					col = s[1:]
					desc = " DESC"
				}
				clauses = append(clauses, col+desc)
			}
			db = db.Order(strings.Join(clauses, ", "))
		}

		if opts.Limit > 0 {
			offset := (opts.Page - 1) * opts.Limit
			db = db.Limit(opts.Limit).Offset(offset)
		}

		return db
	}
}

// sanitizeLikeInput sanitizes input for LIKE operations to prevent injection.
func sanitizeLikeInput(input any) string {
	s := fmt.Sprint(input)

	// Remove or escape LIKE wildcards to prevent unintended matching
	// Note: GORM handles SQL escaping, but we prevent wildcard abuse
	s = strings.ReplaceAll(s, "%", "")
	s = strings.ReplaceAll(s, "_", "")

	// Limit input length to prevent DoS
	if len(s) > MaxInputLength {
		s = s[:MaxInputLength]
	}

	return s
}

// -------------------- Operators --------------------.
func (qb *QueryBuilder) registerDefaultOperators() {
	op := func(expr string) OperatorHandler {
		return func(db *gorm.DB, f string, v any) *gorm.DB {
			return db.Where(f+expr, v)
		}
	}

	qb.operators = map[string]OperatorHandler{
		"eq":  op(" = ?"),
		"neq": op(" != ?"),
		"gt":  op(" > ?"),
		"gte": op(" >= ?"),
		"lt":  op(" < ?"),
		"lte": op(" <= ?"),
		"like": func(db *gorm.DB, f string, v any) *gorm.DB {
			safeValue := sanitizeLikeInput(v)
			return db.Where(f+" LIKE ?", "%"+safeValue+"%")
		},
		"ilike": func(db *gorm.DB, f string, v any) *gorm.DB {
			safeValue := sanitizeLikeInput(v)
			return db.Where(f+" ILIKE ?", "%"+safeValue+"%")
		},
		"in": func(db *gorm.DB, f string, v any) *gorm.DB {
			if arr, ok := v.([]any); ok && len(arr) > MaxArrayElements {
				// Truncate array to prevent DoS
				arr = arr[:MaxArrayElements]
				return db.Where(f+" IN ?", arr)
			}
			return db.Where(f+" IN ?", v)
		},
		"nin": func(db *gorm.DB, f string, v any) *gorm.DB {
			if arr, ok := v.([]any); ok && len(arr) > MaxArrayElements {
				// Truncate array to prevent DoS
				arr = arr[:MaxArrayElements]
				return db.Where(f+" NOT IN ?", arr)
			}
			return db.Where(f+" NOT IN ?", v)
		},
		"null":    func(db *gorm.DB, f string, _ any) *gorm.DB { return db.Where(f + " IS NULL") },
		"notnull": func(db *gorm.DB, f string, _ any) *gorm.DB { return db.Where(f + " IS NOT NULL") },
		"between": func(db *gorm.DB, f string, v any) *gorm.DB {
			if vals, ok := v.([]any); ok && len(vals) == 2 {
				return db.Where(f+" BETWEEN ? AND ?", vals[0], vals[1])
			}
			return db
		},
	}
}

// -------------------- Field Helpers --------------------.

func (qb *QueryBuilder) RegisterStringField(field string, column ...string) {
	qb.RegisterField(field, func(v string) (any, error) { return v, nil }, column...)
}
func (qb *QueryBuilder) RegisterIntField(field string, column ...string) {
	qb.RegisterField(field, func(v string) (any, error) { return strconv.Atoi(v) }, column...)
}
func (qb *QueryBuilder) RegisterBoolField(field string, column ...string) {
	qb.RegisterField(field, func(v string) (any, error) { return strconv.ParseBool(v) }, column...)
}
func (qb *QueryBuilder) RegisterTimeField(field string, column ...string) {
	qb.RegisterField(field, func(v string) (any, error) { return time.Parse(time.RFC3339, v) }, column...)
}
