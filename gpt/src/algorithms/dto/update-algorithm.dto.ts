import { PartialType } from '@nestjs/mapped-types';
import { SaveAlgorithmDto } from './save-algorithm.dto';

export class UpdateAlgorithmDto extends PartialType(SaveAlgorithmDto) {}
