import { Column, CreateDateColumn, Entity, PrimaryColumn } from 'typeorm';

@Entity('password_resets')
export class PasswordReset {
  @PrimaryColumn({ type: 'varchar', length: 255 })
  email!: string;

  @Column({ type: 'varchar', length: 255 })
  token!: string;

  @CreateDateColumn({ type: 'datetime', nullable: true, name: 'created_at' })
  createdAt!: Date | null;
}
