export function calculateLevel(experiencePoint: number): number {
  if (experiencePoint <= 0) return 1;

  const BASE_EXP = 19;
  const RATE_OF_INCREASE = 0.067;

  const numerator = Math.log(
    1 + (experiencePoint * RATE_OF_INCREASE) / BASE_EXP
  );
  const denominator = Math.log(1 + RATE_OF_INCREASE);
  const calculatedLevel = 1 + numerator / denominator;

  return Math.max(1, Math.floor(calculatedLevel));
}
