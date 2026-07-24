export function calculateExperience(level: number): number {
  if (level <= 1) return 0;

  const BASE_EXP = 19;
  const RATE_OF_INCREASE = 0.067;

  const factor = BASE_EXP / RATE_OF_INCREASE;
  const power = Math.pow(1 + RATE_OF_INCREASE, level - 1);
  const requiredExperience = factor * (power - 1);

  return Math.ceil(requiredExperience);
}
