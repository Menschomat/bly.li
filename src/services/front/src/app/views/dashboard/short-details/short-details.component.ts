import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import {
  catchError,
  filter,
  map,
  Observable,
  of,
  shareReplay,
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
  host: { class: 'flex-1 flex  justify-center' },
  template: `
    <main
      class="w-full max-w-6xl px-4 py-6 sm:px-6 lg:px-8 flex flex-col gap-4 overflow-auto"
    >
      @if (clickHistory$ | async; as result) { @if (result.status === 'loading')
      {
      <div class="flex-1 flex justify-center items-center">
        <h3 class="text-xl"><b>Loading…</b></h3>
      </div>
      } @else if (result.status === 'error') {
      <div class="flex-1 flex justify-center items-center">
        <h3 class="text-xl">
          <b class="text-red-600 dark:text-red-400">Error:</b>
          {{ result.message }}
        </h3>
      </div>
      } @else if (result.status === 'success') { @if (result.data.length === 0)
      {
      <div class="flex-1 flex justify-center items-center">
        <h3 class="text-xl">
          <b>No data yet.</b> Looks like your short URL hasn't been clicked.
          <b>Share it to get started!</b>
        </h3>
      </div>
      } @else {
      <header>
        <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
          <h1 class="text-3xl font-bold">Dashboard</h1>
          <h3 class="py-2 text-xl">
            Details of Short <b>{{ 'TODO' }}</b>
          </h3>
        </div>
      </header>
      <div
        class="
        mx-auto max-w-7xl px-4 py-6 
        sm:px-6 lg:px-8
        flex flex-col gap-4 
        w-full
        backdrop-blur-md 
        rounded-3xl 
        border border-gray-100 
        shadow-sm 
        dark:border-gray-900"
      >
        <app-line-chart
          class="max-h-[30rem] min-h-[20rem]"
          [chartTitle]="'Click histogram'"
          [data$]="chartClickHistory$"
        ></app-line-chart>
      </div>

     <!--<ul>
        @for (item of result.data; track item.Timestamp) {
        <li>{{ item.Timestamp }} – {{ item.Count }}</li>
        }
      </ul>-->
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

  ngOnInit(): void {
    // 1️⃣ Build a shared, hot stream of ClickHistoryState
    const sharedClicks$ = this.route.paramMap.pipe(
      switchMap((params) => {
        console.log('UPDATE');
        const id = params.get('short') ?? '';
        return this.dashboardService.clickHistoryByShort(id).pipe(
          // successful response
          map(
            (data) =>
              ({
                status: 'success' as const,
                data,
              } as { status: 'success'; data: ShortClickCount[] })
          ),
          // handle errors
          catchError((error) =>
            of({
              status: 'error' as const,
              code: error.status,
              message:
                error.status === 403
                  ? 'Short does not exist, or you are not authorized to view it.'
                  : error.message ?? 'An error occurred.',
            } as { status: 'error'; code: number; message: string })
          ),
          // initial loading state
          startWith({ status: 'loading' } as const)
        );
      }),
      // share the last ClickHistoryState with ALL subscribers
      shareReplay({ bufferSize: 1, refCount: true })
    );

    // 2️⃣ Expose for template async pipe (loading/error UI)
    this.clickHistory$ = sharedClicks$;

    // 3️⃣ Derive DataPoint[] for your charts (y is definitely a number)
    this.chartClickHistory$ = sharedClicks$.pipe(
      filter(
        (h): h is { status: 'success'; data: ShortClickCount[] } =>
          h.status === 'success'
      ),
      map((h) =>
        h.data.map(
          (a) =>
            ({
              x: new Date(a.Timestamp ?? 0).getTime(),
              y: a.Count ?? 0,
            } as DataPoint)
        )
      )
    );
  }
}

export type ClickHistoryState =
  | { status: 'loading' }
  | { status: 'error'; code: number; message: string }
  | { status: 'success'; data: ShortClickCount[] };
