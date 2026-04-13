import {
  Column,
  CreateDateColumn,
  DeleteDateColumn,
  Entity,
  Index,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';
import {
  ADMIN_USER_TYPES,
  ADMIN_VERIFIED_VALUES,
} from '../../common/constants';

@Entity('admins')
export class Admin {
  @PrimaryGeneratedColumn()
  id!: number;

  @Column({ type: 'varchar', length: 255 })
  firstname!: string;

  @Column({ type: 'varchar', length: 255 })
  lastname!: string;

  @Column({ type: 'varchar', length: 32, nullable: true })
  phone!: string | null;

  @Index({ unique: true })
  @Column({ type: 'varchar', length: 255 })
  email!: string;

  @Column({ type: 'simple-enum', enum: ADMIN_VERIFIED_VALUES, default: 'N' })
  verified!: (typeof ADMIN_VERIFIED_VALUES)[number];

  @Column({ type: 'simple-enum', enum: ADMIN_USER_TYPES })
  userType!: (typeof ADMIN_USER_TYPES)[number];

  @Column({ type: 'varchar', length: 255 })
  password!: string;

  @Column({ type: 'text', nullable: true })
  accessToken!: string | null;

  @Column({ type: 'varchar', length: 255, nullable: true, name: 'remember_token' })
  rememberToken!: string | null;

  @DeleteDateColumn({ type: 'datetime', nullable: true, name: 'deleted_at' })
  deletedAt!: Date | null;

  @CreateDateColumn({ type: 'datetime', name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ type: 'datetime', name: 'updated_at' })
  updatedAt!: Date;
}
