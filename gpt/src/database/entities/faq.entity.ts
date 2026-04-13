import {
  Column,
  CreateDateColumn,
  DeleteDateColumn,
  Entity,
  Index,
  JoinColumn,
  OneToOne,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';
import { FaqGroup } from './faq-group.entity';

@Entity('faqs')
export class Faq {
  @PrimaryGeneratedColumn()
  id!: number;

  @Index({ unique: true })
  @Column({ type: 'integer', name: 'faq_group_id' })
  faqGroupId!: number;

  @Column({ type: 'simple-json', nullable: true })
  content!: Record<string, unknown> | null;

  @DeleteDateColumn({ type: 'datetime', nullable: true, name: 'deleted_at' })
  deletedAt!: Date | null;

  @CreateDateColumn({ type: 'datetime', name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ type: 'datetime', name: 'updated_at' })
  updatedAt!: Date;

  @OneToOne(() => FaqGroup, (group) => group.faq, { nullable: false })
  @JoinColumn({ name: 'faq_group_id' })
  group?: FaqGroup;
}
