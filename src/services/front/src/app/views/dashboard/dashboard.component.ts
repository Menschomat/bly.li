import { Component } from '@angular/core';
import { DasherService, ShortURL } from '../../api';
import { Observable } from 'rxjs';
import { CommonModule } from '@angular/common';
import { NumberToWordsPipe } from '../../pipes/number-to-words.pipe';
import { ShortTableComponent } from './short-table/short-table.component';
import { DashboardService } from '../../services/dashboard.service';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule, NumberToWordsPipe, ShortTableComponent],
  host: { class: 'flex-1 ' },
  template: `
    <header>
      <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <h1 class="text-3xl font-bold  ">Dashboard</h1>
        <h3 class="py-2">
          Welcome back Mr. Example. Currently you have
          <b>{{ ($allShorts | async)?.length ?? 0 | numberToWords }}</b> shorts.
        </h3>
      </div>
    </header>
    <main
      class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8 flex flex-col gap-4"
    >
      <app-short-table></app-short-table>
    </main>
  `,
  styles: ``,
})
export class DashboardComponent {
  public $allShorts: Observable<ShortURL[]>;
  constructor(private dasherService: DashboardService) {
    this.$allShorts = this.dasherService.$allShorts;
  }
}
