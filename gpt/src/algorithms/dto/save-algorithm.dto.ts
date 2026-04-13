import {
  IsIn,
  IsInt,
  IsOptional,
  IsString,
  MaxLength,
} from 'class-validator';
import {
  ACTIVE_FLAG_VALUES,
  ALGORITHM_ACTION_TYPES,
  PROGESTOGEN_DIRECTIONS,
  PROGESTOGEN_POSSIBLE,
} from '../../common/constants';

export class SaveAlgorithmDto {
  @IsString()
  text!: string;

  @IsOptional()
  @IsIn(ALGORITHM_ACTION_TYPES)
  actionType?: (typeof ALGORITHM_ACTION_TYPES)[number] | null;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  positive?: string | null;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  negative?: string | null;

  @IsOptional()
  @IsInt()
  onPositive?: number | null;

  @IsOptional()
  @IsInt()
  onNegative?: number | null;

  @IsOptional()
  @IsInt()
  nextMove?: number | null;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  tempPlan?: string | null;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  tempPlanDirP?: string | null;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  tempPlanDirN?: string | null;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  conditionalFactor?: string | null;

  @IsOptional()
  @IsString()
  @MaxLength(32)
  conditionalOperator?: string | null;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  conditionalValue?: string | null;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  stateValue?: string | null;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  label?: string | null;

  @IsOptional()
  @IsIn(PROGESTOGEN_POSSIBLE)
  progestogenPossible?: (typeof PROGESTOGEN_POSSIBLE)[number] | null;

  @IsOptional()
  @IsIn(PROGESTOGEN_DIRECTIONS)
  progestogenPossibleDir?: (typeof PROGESTOGEN_DIRECTIONS)[number] | null;

  @IsInt()
  delay!: number;

  @IsOptional()
  @IsInt()
  series?: number | null;

  @IsOptional()
  @IsIn(ACTIVE_FLAG_VALUES)
  active?: (typeof ACTIVE_FLAG_VALUES)[number];
}
