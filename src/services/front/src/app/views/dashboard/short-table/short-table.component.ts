import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  NO_ERRORS_SCHEMA,
} from '@angular/core';
import { DasherService, ShortURL } from '../../../api';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { CommonModule } from '@angular/common';
import { ShortTableRowComponent } from './short-table-row/short-table-row.component';
import { DashboardService } from '../../../services/dashboard.service';
@Component({
  selector: 'app-short-table',
  imports: [CommonModule, ShortTableRowComponent],
  schemas: [CUSTOM_ELEMENTS_SCHEMA, NO_ERRORS_SCHEMA],
  host: { class: 'flex-1 flex flex-col gap-4' },
  template: `
    @for (item of $allShorts | async; track item.Short) {
    <app-short-table-row (delete)="delete(item)"
      ><row-title>{{ item.Short }}</row-title>
      <row-url>{{ item.URL }}</row-url>
      <row-count>{{ '10k' }} Clicks</row-count></app-short-table-row
    >
    }
  `,
  styles: ``,
})
export class ShortTableComponent {
  public $allShorts: Observable<ShortURL[]>;

  constructor(private dasherService: DashboardService) {
    this.$allShorts = dasherService.$allShorts;
  }
  public delete(short: ShortURL) {
    this.dasherService.delete(short);
  }
}
