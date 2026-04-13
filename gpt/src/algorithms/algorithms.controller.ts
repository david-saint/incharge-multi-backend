import { Body, Controller, Get, Param, ParseIntPipe, Post, Put, Query, UseGuards } from '@nestjs/common';
import { VerifiedAdminGuard } from '../common/guards/verified-admin.guard';
import { AlgorithmsService } from './algorithms.service';
import { SaveAlgorithmDto } from './dto/save-algorithm.dto';
import { UpdateAlgorithmDto } from './dto/update-algorithm.dto';

@Controller('algo')
export class AlgorithmsController {
  constructor(private readonly algorithmsService: AlgorithmsService) {}

  @Get()
  listAlgorithms(@Query('active') active?: string) {
    return this.algorithmsService.list(active);
  }

  @Post()
  @UseGuards(VerifiedAdminGuard)
  createAlgorithm(@Body() payload: SaveAlgorithmDto) {
    return this.algorithmsService.create(payload);
  }

  @Put(':id')
  @UseGuards(VerifiedAdminGuard)
  updateAlgorithm(
    @Param('id', ParseIntPipe) id: number,
    @Body() payload: UpdateAlgorithmDto,
  ) {
    return this.algorithmsService.update(id, payload);
  }
}
