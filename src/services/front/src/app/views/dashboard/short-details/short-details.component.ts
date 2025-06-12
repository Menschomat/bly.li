import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import {
  catchError,
  filter,
  map,
  Observable,
  of,
  startWith,
  switchMap,
} from 'rxjs';
import { DashboardService } from '../../../services/dashboard.service';
import { ShortClickCount } from '../../../api';
import { CommonModule } from '@angular/common';
import {
  DataPoint,
  LineChartComponent,
} from '../../../components/graphs/line-chart/line-chart.component';

@Component({
  selector: 'app-short-details',
  imports: [CommonModule, LineChartComponent],
  host: { class: 'flex-1 flex w-full  justify-center' },
  template: `
    <main
      class="w-full max-w-6xl px-4 py-6 sm:px-6 lg:px-8 flex flex-col gap-4 overflow-auto"
    >
      @if (clickHistory$ | async; as result) { @if (result.status === 'loading')
      {
      <div class="flex-1 flex items-center">
        <h3 class="text-xl"><b>Loading…</b></h3>
      </div>
      } @else if (result.status === 'error') {
      <div class="flex-1 flex items-center">
        <h3 class="text-xl">
          <b class="text-red-600 dark:text-red-400">Error:</b>
          {{ result.message }}
        </h3>
      </div>
      } @else if (result.status === 'success') { @if (result.data.length === 0)
      {
      <div class="flex-1 flex items-center">
        <h3 class="text-xl">
          <b>No data yet.</b> Looks like your short URL hasn't been clicked.
          <b>Share it to get started!</b>
        </h3>
      </div>
      } @else {
      <app-line-chart
      class="max-h-[30rem] min-h-[20rem]"
        [chartTitle]="'Click histogram'"
        [data$]="chartClickHistory$"
      ></app-line-chart>
      <!--<app-line-chart
      class="max-h-[30rem] min-h-[20rem]"
        [chartTitle]="'Click histogram'"
        [data$]="chartClickHistory$"
      ></app-line-chart>
            <app-line-chart
      class="max-h-[30rem] min-h-[20rem]"
        [chartTitle]="'Click histogram'"
        [data$]="chartClickHistory$"
      ></app-line-chart>
            <app-line-chart
      class="max-h-[30rem] min-h-[20rem]"
        [chartTitle]="'Click histogram'"
        [data$]="chartClickHistory$"
      ></app-line-chart>-->
      <ul>
        @for (item of result.data; track item.Timestamp) {
        <li>{{ item.Timestamp }} – {{ item.Count }}</li>
        }
      </ul>
      } } }
    </main>
  `,
  styles: ``,
})
export class ShortDetailsComponent implements OnInit {
  clickHistory$!: Observable<ClickHistoryState>;
  chartClickHistory$!: Observable<DataPoint[]>;

  constructor(
    private route: ActivatedRoute,
    private dashboardService: DashboardService
  ) {}

  ngOnInit() {
    this.clickHistory$ = this.route.paramMap.pipe(
      switchMap((params) => {
        const id = params.get('short');
        return this.dashboardService.clickHistoryByShort(id ?? '').pipe(
          map(
            (data: ShortClickCount[]) =>
              ({ status: 'success', data: data ?? [] } as ClickHistoryState)
          ),
          catchError((error) => {
            if (error.status === 403) {
              return of({
                status: 'error',
                code: 403,
                message:
                  'Short does not exists, or you are not authorized to view it.',
              } as ClickHistoryState);
            }
            // Handle other errors
            return of({
              status: 'error',
              code: error.status,
              message: error.message || 'An error occurred.',
            } as ClickHistoryState);
          }),
          startWith({ status: 'loading' } as ClickHistoryState)
        );
      })
    );
    this.chartClickHistory$ = this.clickHistory$.pipe(
      filter((hist) => hist.status === 'success'),
      map((hist) => {
        return hist.data.map((a) => {
          return {
            x: new Date(a.Timestamp ?? 0).getTime(),
            y: a.Count,
          } as DataPoint;
        });
      })
    );
  }
}

export type ClickHistoryState =
  | { status: 'loading' }
  | { status: 'error'; code: number; message: string }
  | { status: 'success'; data: ShortClickCount[] };
