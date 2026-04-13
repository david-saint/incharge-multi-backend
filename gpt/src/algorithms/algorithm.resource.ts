import { Algorithm } from '../database/entities/algorithm.entity';

export function buildAlgorithmResource(algorithm: Algorithm) {
  return {
    id: algorithm.id,
    text: algorithm.text,
    actionType: algorithm.actionType,
    positive: algorithm.positive,
    negative: algorithm.negative,
    onPositive: algorithm.onPositive,
    onNegative: algorithm.onNegative,
    nextMove: algorithm.nextMove,
    tempPlan: algorithm.tempPlan,
    tempPlanDirP: algorithm.tempPlanDirP,
    tempPlanDirN: algorithm.tempPlanDirN,
    conditionalFactor: algorithm.conditionalFactor,
    conditionalOperator: algorithm.conditionalOperator,
    conditionalValue: algorithm.conditionalValue,
    stateValue: algorithm.stateValue,
    label: algorithm.label,
    progestogenPossible: algorithm.progestogenPossible,
    progestogenPossibleDir: algorithm.progestogenPossibleDir,
    delay: algorithm.delay,
    series: algorithm.series,
    active: algorithm.active,
    created_at: algorithm.createdAt,
    updated_at: algorithm.updatedAt,
  };
}
