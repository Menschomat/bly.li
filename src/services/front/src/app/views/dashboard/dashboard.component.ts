import { Component } from '@angular/core';
import { DasherService, ShortURL } from '../../api';
import { Observable } from 'rxjs';
import { CommonModule } from '@angular/common';
import { NumberToWordsPipe } from '../../pipes/number-to-words.pipe';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule, NumberToWordsPipe],
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
      @for (item of $allShorts | async; track item.Short) {
      <div class="flex flex-col gap-4">
        <div>
          {{ item.Short }}
        </div>
      </div>
      }
    </main>
  `,
  styles: ``,
})
export class DashboardComponent {
  public $allShorts: Observable<ShortURL[]>;
  constructor(private dasherService: DasherService) {
    this.$allShorts = this.dasherService.dasherShortAllGet();
  }
}
