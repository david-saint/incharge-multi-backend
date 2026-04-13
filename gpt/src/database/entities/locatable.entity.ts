import {
  Column,
  CreateDateColumn,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';
import { Clinic } from './clinic.entity';
import { Location } from './location.entity';

@Entity('locatables')
export class Locatable {
  @PrimaryGeneratedColumn()
  id!: number;

  @Column({ type: 'integer', name: 'location_id' })
  locationId!: number;

  @Column({ type: 'integer', name: 'locatable_id' })
  locatableId!: number;

  @Column({ type: 'varchar', length: 255, name: 'locatable_type' })
  locatableType!: string;

  @CreateDateColumn({ type: 'datetime', name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ type: 'datetime', name: 'updated_at' })
  updatedAt!: Date;

  @ManyToOne(() => Location, (location) => location.locatables, { nullable: false })
  @JoinColumn({ name: 'location_id' })
  location?: Location;

  @ManyToOne(() => Clinic, (clinic) => clinic.locatables, {
    nullable: true,
    createForeignKeyConstraints: false,
  })
  @JoinColumn({ name: 'locatable_id', referencedColumnName: 'id' })
  clinic?: Clinic;
}
