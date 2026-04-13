import { Controller, Get, Param, ParseIntPipe } from '@nestjs/common';
import { GlobalService } from './global.service';

@Controller('api/v1/global')
export class GlobalController {
  constructor(private readonly globalService: GlobalService) {}

  @Get()
  hello() {
    return this.globalService.hello();
  }

  @Get('contraception-reasons')
  listContraceptionReasons() {
    return this.globalService.listContraceptionReasons();
  }

  @Get('contraception-reasons/:id')
  getContraceptionReason(@Param('id', ParseIntPipe) id: number) {
    return this.globalService.getContraceptionReason(id);
  }

  @Get('education-levels')
  listEducationLevels() {
    return this.globalService.listEducationLevels();
  }

  @Get('faq-groups')
  listFaqGroups() {
    return this.globalService.listFaqGroups();
  }

  @Get('faq-groups/:id')
  getFaqContent(@Param('id', ParseIntPipe) id: number) {
    return this.globalService.getFaqGroupContent(id);
  }
}
