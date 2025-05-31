import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  NO_ERRORS_SCHEMA,
} from '@angular/core';
import { ShortURL } from '../../../api';
import { Observable } from 'rxjs';
import { CommonModule } from '@angular/common';
import { ShortTableRowComponent } from './short-table-row/short-table-row.component';
import { DashboardService } from '../../../services/dashboard.service';
import { ShortNumberPipe } from '../../../pipes/short-number.pipe';
import { ConfigService } from '../../../services/config.service';
@Component({
  selector: 'app-short-table',
  imports: [CommonModule, ShortNumberPipe, ShortTableRowComponent],
  schemas: [CUSTOM_ELEMENTS_SCHEMA, NO_ERRORS_SCHEMA],
  host: { class: 'flex-1 flex flex-col gap-4' },
  template: `
    @for (item of $allShorts | async; track item.Short) {
    <app-short-table-row (copy)="copyItem(item)" (delete)="deleteItem(item)"
      ><row-title>{{ item.Short }}</row-title>
      <row-url>{{ item.URL }}</row-url>
      <row-count
        >{{ item.Count ?? 0 | shortNumber }} Clicks</row-count
      ></app-short-table-row
    >
    }
  `,
  styles: ``,
})
export class ShortTableComponent {
  public $allShorts: Observable<ShortURL[]>;

  constructor(
    private dasherService: DashboardService,
    private config: ConfigService
  ) {
    this.$allShorts = dasherService.$allShorts;
  }
  public async deleteItem(short: ShortURL) {
    await this.dasherService.delete(short);
  }
  public async copyItem(short: ShortURL) {
    try {
      await navigator.clipboard.writeText(
        `${window.location.origin}${this.config.getConfig().blowUpPath}/${
          short.Short
        }`
      );
    } catch (err) {
      console.error('Failed to copy: ', err);
    }
  }
}
