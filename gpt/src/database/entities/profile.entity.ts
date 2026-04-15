import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToOne,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';
import {
  GENDERS,
  MARITAL_STATUSES,
  RELIGIONS,
  RELIGION_SECTS,
} from '../../common/constants';
import { ContraceptionReason } from './contraception-reason.entity';
import { EducationLevel } from './education-level.entity';
import { User } from './user.entity';

@Entity('profiles')
export class Profile {
  @PrimaryGeneratedColumn()
  id!: number;

  @Index({ unique: true })
  @Column({ type: 'integer', name: 'user_id' })
  userId!: number;

  @OneToOne(() => User, (user) => user.profile, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'user_id' })
  user!: User;

  @Column({ type: 'integer', unsigned: true, default: 0 })
  age!: number;

  @Column({ type: 'simple-enum', enum: GENDERS })
  gender!: (typeof GENDERS)[number];

  @Column({ type: 'datetime', name: 'date_of_birth' })
  dateOfBirth!: Date;

  @Column({ type: 'text' })
  address!: string;

  @Column({ type: 'decimal', precision: 10, scale: 7, nullable: true })
  latitude!: string | null;

  @Column({ type: 'decimal', precision: 10, scale: 7, nullable: true })
  longitude!: string | null;

  @Column({
    type: 'simple-enum',
    enum: MARITAL_STATUSES,
    name: 'marital_status',
    default: 'SINGLE',
  })
  maritalStatus!: (typeof MARITAL_STATUSES)[number];

  @Column({ type: 'integer', unsigned: true, nullable: true })
  height!: number | null;

  @Column({ type: 'decimal', precision: 10, scale: 2, nullable: true })
  weight!: string | null;

  @Column({ type: 'integer', nullable: true, name: 'education_level_id' })
  educationLevelId!: number | null;

  @ManyToOne(() => EducationLevel, (educationLevel) => educationLevel.profiles, {
    nullable: true,
  })
  @JoinColumn({ name: 'education_level_id' })
  educationLevel?: EducationLevel | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  occupation!: string | null;

  @Column({ type: 'integer', unsigned: true, nullable: true, name: 'number_of_children' })
  numberOfChildren!: number | null;

  @Column({ type: 'integer', nullable: true, name: 'contraception_reason_id' })
  contraceptionReasonId!: number | null;

  @ManyToOne(() => ContraceptionReason, (reason) => reason.profiles, { nullable: true })
  @JoinColumn({ name: 'contraception_reason_id' })
  contraceptionReason?: ContraceptionReason | null;

  @Column({ type: 'boolean', name: 'sexually_active', default: false })
  sexuallyActive!: boolean;

  @Column({ type: 'boolean', name: 'pregnancy_status', default: false })
  pregnancyStatus!: boolean;

  @Column({ type: 'simple-enum', enum: RELIGIONS, default: 'OTHER' })
  religion!: (typeof RELIGIONS)[number];

  @Column({
    type: 'simple-enum',
    enum: RELIGION_SECTS,
    nullable: true,
    name: 'religion_sect',
  })
  religionSect!: (typeof RELIGION_SECTS)[number] | null;

  @Column({ type: 'simple-json', nullable: true })
  meta!: Record<string, unknown> | null;

  @CreateDateColumn({ type: 'datetime', name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ type: 'datetime', name: 'updated_at' })
  updatedAt!: Date;
}
