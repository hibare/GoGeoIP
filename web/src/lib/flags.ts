import { hasFlag } from "country-flag-icons";
import getUnicodeFlagIcon from "country-flag-icons/unicode";

export { hasFlag };

export function getCountryFlag(countryCode: string): string | null {
  if (!countryCode || !hasFlag(countryCode)) return null;
  return getUnicodeFlagIcon(countryCode.toUpperCase());
}
