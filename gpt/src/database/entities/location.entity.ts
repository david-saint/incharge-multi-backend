import {
  Column,
  CreateDateColumn,
  DeleteDateColumn,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';
import { Country } from './country.entity';
import { Locatable } from './locatable.entity';
import { State } from './state.entity';

@Entity('locations')
export class Location {
  @PrimaryGeneratedColumn()
  id!: number;

  @Column({ type: 'varchar', length: 255 })
  name!: string;

  @Column({ type: 'integer', name: 'state_id' })
  stateId!: number;

  @Column({ type: 'integer', name: 'country_id' })
  countryId!: number;

  @Column({ type: 'decimal', precision: 10, scale: 7, nullable: true })
  latitude!: string | null;

  @Column({ type: 'decimal', precision: 10, scale: 7, nullable: true })
  longitude!: string | null;

  @DeleteDateColumn({ type: 'datetime', nullable: true, name: 'deleted_at' })
  deletedAt!: Date | null;

  @Column({ type: 'simple-json', nullable: true })
  meta!: Record<string, unknown> | null;

  @CreateDateColumn({ type: 'datetime', name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ type: 'datetime', name: 'updated_at' })
  updatedAt!: Date;

  @ManyToOne(() => State, (state) => state.locations, { nullable: false })
  @JoinColumn({ name: 'state_id' })
  state?: State;

  @ManyToOne(() => Country, (country) => country.locations, { nullable: false })
  @JoinColumn({ name: 'country_id' })
  country?: Country;

  @OneToMany(() => Locatable, (locatable) => locatable.location)
  locatables?: Locatable[];
}
