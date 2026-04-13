import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Algorithm } from '../database/entities/algorithm.entity';
import { AlgorithmsController } from './algorithms.controller';
import { AlgorithmsService } from './algorithms.service';

@Module({
  imports: [TypeOrmModule.forFeature([Algorithm])],
  controllers: [AlgorithmsController],
  providers: [AlgorithmsService],
  exports: [AlgorithmsService],
})
export class AlgorithmsModule {}
