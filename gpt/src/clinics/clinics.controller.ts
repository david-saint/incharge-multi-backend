import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  ParseIntPipe,
  Post,
  Put,
  Query,
  UseGuards,
} from '@nestjs/common';
import { AdminSessionGuard } from '../common/guards/admin-session.guard';
import { VerifiedAdminGuard } from '../common/guards/verified-admin.guard';
import { ClinicListQueryDto } from './dto/clinic-list-query.dto';
import { SaveClinicDto } from './dto/save-clinic.dto';
import { ClinicsService } from './clinics.service';

@Controller()
export class ClinicsController {
  constructor(private readonly clinicsService: ClinicsService) {}

  @Get('api/v1/user/clinics')
  listClinics(@Query() query: ClinicListQueryDto) {
    return this.clinicsService.list(query, '/api/v1/user/clinics');
  }

  @Get('api/v1/user/clinics/getClinics')
  @UseGuards(AdminSessionGuard)
  listSessionClinics() {
    return this.clinicsService.listSimple(false, '/api/v1/user/clinics/getClinics');
  }

  @Get('api/v1/user/clinics/deletedClinics')
  @UseGuards(AdminSessionGuard)
  listDeletedSessionClinics() {
    return this.clinicsService.listSimple(true, '/api/v1/user/clinics/deletedClinics');
  }

  @Post('api/v1/user/clinics/addClinic')
  @UseGuards(AdminSessionGuard)
  createSessionClinic(@Body() payload: SaveClinicDto) {
    return this.clinicsService.create(payload);
  }

  @Put('api/v1/user/clinics/update/:clinic_id')
  @UseGuards(AdminSessionGuard)
  updateSessionClinic(
    @Param('clinic_id', ParseIntPipe) clinicId: number,
    @Body() payload: SaveClinicDto,
  ) {
    return this.clinicsService.update(clinicId, payload);
  }

  @Put('api/v1/user/clinics/revertDelete/:clinic_id')
  @UseGuards(AdminSessionGuard)
  restoreSessionClinic(@Param('clinic_id', ParseIntPipe) clinicId: number) {
    return this.clinicsService.restore(clinicId);
  }

  @Delete('api/v1/user/clinics/deleteClinic/:clinic_id')
  @UseGuards(AdminSessionGuard)
  deleteSessionClinic(@Param('clinic_id', ParseIntPipe) clinicId: number) {
    return this.clinicsService.softDelete(clinicId);
  }

  @Get('getClinics')
  @UseGuards(VerifiedAdminGuard)
  getClinics() {
    return this.clinicsService.listSimple(false);
  }

  @Get('getDeletedClinics')
  @UseGuards(VerifiedAdminGuard)
  getDeletedClinics() {
    return this.clinicsService.listSimple(true);
  }

  @Post('addClinic')
  @UseGuards(VerifiedAdminGuard)
  addClinic(@Body() payload: SaveClinicDto) {
    return this.clinicsService.create(payload);
  }

  @Put('updateClinic/:clinic_id')
  @UseGuards(VerifiedAdminGuard)
  updateClinic(
    @Param('clinic_id', ParseIntPipe) clinicId: number,
    @Body() payload: SaveClinicDto,
  ) {
    return this.clinicsService.update(clinicId, payload);
  }

  @Delete('deleteClinic/:clinic_id')
  @UseGuards(VerifiedAdminGuard)
  deleteClinic(@Param('clinic_id', ParseIntPipe) clinicId: number) {
    return this.clinicsService.softDelete(clinicId);
  }

  @Put('revertDeletedClinic/:clinic_id')
  @UseGuards(VerifiedAdminGuard)
  revertDeletedClinic(@Param('clinic_id', ParseIntPipe) clinicId: number) {
    return this.clinicsService.restore(clinicId);
  }
}
