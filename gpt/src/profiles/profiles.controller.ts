import { Body, Controller, Get, HttpCode, Post, Query, UseGuards } from '@nestjs/common';
import { CurrentUser } from '../common/decorators/current-user.decorator';
import { JwtAuthGuard } from '../common/guards/jwt-auth.guard';
import { parseCsvList } from '../common/utils/query.util';
import { User } from '../database/entities/user.entity';
import { SaveProfileDto } from './dto/save-profile.dto';
import { StorePlanDto } from './dto/store-plan.dto';
import { ProfilesService } from './profiles.service';

@Controller('api/v1/user/profile')
@UseGuards(JwtAuthGuard)
export class ProfilesController {
  constructor(private readonly profilesService: ProfilesService) {}

  @Post()
  async saveProfile(@CurrentUser() user: User, @Body() payload: SaveProfileDto) {
    return this.profilesService.upsertProfile(user.id, payload);
  }

  @Get()
  async getProfile(@CurrentUser() user: User, @Query('with') withQuery?: string) {
    return this.profilesService.getProfile(user.id, parseCsvList(withQuery));
  }

  @Post('algorithm')
  @HttpCode(200)
  async storeAlgorithmPlan(@CurrentUser() user: User, @Body() payload: StorePlanDto) {
    return this.profilesService.storePlan(user.id, payload.plan);
  }
}
