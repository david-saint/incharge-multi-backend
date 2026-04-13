import {
  Column,
  CreateDateColumn,
  DeleteDateColumn,
  Entity,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';
import {
  ACTIVE_FLAG_VALUES,
  ALGORITHM_ACTION_TYPES,
  PROGESTOGEN_DIRECTIONS,
  PROGESTOGEN_POSSIBLE,
} from '../../common/constants';

@Entity('algorithms')
export class Algorithm {
  @PrimaryGeneratedColumn()
  id!: number;

  @Column({ type: 'text' })
  text!: string;

  @Column({ type: 'simple-enum', enum: ALGORITHM_ACTION_TYPES, nullable: true })
  actionType!: (typeof ALGORITHM_ACTION_TYPES)[number] | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  positive!: string | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  negative!: string | null;

  @Column({ type: 'integer', nullable: true })
  onPositive!: number | null;

  @Column({ type: 'integer', nullable: true })
  onNegative!: number | null;

  @Column({ type: 'integer', nullable: true })
  nextMove!: number | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  tempPlan!: string | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  tempPlanDirP!: string | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  tempPlanDirN!: string | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  conditionalFactor!: string | null;

  @Column({ type: 'varchar', length: 32, nullable: true })
  conditionalOperator!: string | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  conditionalValue!: string | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  stateValue!: string | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  label!: string | null;

  @Column({ type: 'simple-enum', enum: PROGESTOGEN_POSSIBLE, nullable: true })
  progestogenPossible!: (typeof PROGESTOGEN_POSSIBLE)[number] | null;

  @Column({ type: 'simple-enum', enum: PROGESTOGEN_DIRECTIONS, nullable: true })
  progestogenPossibleDir!: (typeof PROGESTOGEN_DIRECTIONS)[number] | null;

  @Column({ type: 'integer', default: 0 })
  delay!: number;

  @Column({ type: 'integer', nullable: true })
  series!: number | null;

  @Column({ type: 'simple-enum', enum: ACTIVE_FLAG_VALUES, default: 'N' })
  active!: (typeof ACTIVE_FLAG_VALUES)[number];

  @DeleteDateColumn({ type: 'datetime', nullable: true, name: 'deleted_at' })
  deletedAt!: Date | null;

  @CreateDateColumn({ type: 'datetime', name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ type: 'datetime', name: 'updated_at' })
  updatedAt!: Date;
}
