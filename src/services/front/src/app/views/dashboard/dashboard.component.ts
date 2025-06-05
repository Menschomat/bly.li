import { Component } from '@angular/core';
import { ShortURL } from '../../api';
import { filter, map, Observable, tap } from 'rxjs';
import { CommonModule } from '@angular/common';
import { NumberToWordsPipe } from '../../pipes/number-to-words.pipe';
import { ShortTableComponent } from './short-table/short-table.component';
import { DashboardService } from '../../services/dashboard.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule, NumberToWordsPipe, ShortTableComponent],
  host: { class: 'flex-1 ' },
  template: `
    <header>
      <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <h1 class="text-3xl font-bold  ">Dashboard</h1>
        <h3 class="text-xl py-2" *ngIf="$curFullName | async as fullname">
          Welcome back <b>{{ fullname }}</b
          >. Currently you have
          <b>{{ ($allShorts | async)?.length ?? 0 | numberToWords }}</b>
          {{ (($allShorts | async)?.length ?? 0) > 1 ? 'shorts' : 'short' }}.
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
  public $curFullName: Observable<string | undefined>;
  constructor(
    private readonly dasherService: DashboardService,
    private readonly auth: AuthService
  ) {
    this.$allShorts = this.dasherService.$allShorts;

    this.$curFullName = auth.currentUser$.pipe(
      filter((a) => a !== null),
      tap((a) => console.debug(a)),
      map((a) => a['name'])
    );
  }
}
