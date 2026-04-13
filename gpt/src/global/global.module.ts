import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ContraceptionReason } from '../database/entities/contraception-reason.entity';
import { EducationLevel } from '../database/entities/education-level.entity';
import { Faq } from '../database/entities/faq.entity';
import { FaqGroup } from '../database/entities/faq-group.entity';
import { GlobalController } from './global.controller';
import { GlobalService } from './global.service';

@Module({
  imports: [TypeOrmModule.forFeature([ContraceptionReason, EducationLevel, FaqGroup, Faq])],
  controllers: [GlobalController],
  providers: [GlobalService],
  exports: [GlobalService],
})
export class GlobalApiModule {}
