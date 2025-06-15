import { Component } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import {
  Observable,
  catchError,
  filter,
  map,
  of,
  shareReplay,
  startWith,
  switchMap,
} from 'rxjs';
import { DashboardService } from '../../../services/dashboard.service';
import { ShortClickCount, ShortURL } from '../../../api';
import { CommonModule } from '@angular/common';
import {
  DataPoint,
  LineChartComponent,
} from '../../../components/graphs/line-chart/line-chart.component';
import { NumberToWordsPipe } from '../../../pipes/number-to-words.pipe';

type LoadState<T> =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'success'; data: T };

@Component({
  selector: 'app-short-details',
  standalone: true,
  imports: [CommonModule, LineChartComponent, NumberToWordsPipe],
  host: { class: 'flex-1 flex justify-center' },
  template: `
    <main
      class="w-full max-w-6xl px-4 py-6 sm:px-6 lg:px-8 flex flex-col gap-4 overflow-auto"
    >
      @if (clickHistory$ | async; as history) { @switch (history.status) { @case
      ('loading') {
      <div class="flex-1 flex justify-center items-center">
        <h3 class="text-xl"><b>Loading…</b></h3>
      </div>
      } @case ('error') {
      <div class="flex-1 flex justify-center items-center">
        <h3 class="text-xl">
          <b class="text-red-600 dark:text-red-400">Error:</b>
          {{ history.message }}
        </h3>
      </div>
      } @case ('success') { @if (history.data.length === 0) {
      <div class="flex-1 flex justify-center items-center">
        <h3 class="text-xl">
          <b>No data yet.</b> Looks like your short URL hasn't been clicked.
          <b>Share it to get started!</b>
        </h3>
      </div>
      } @else {
      <header>
        @if (shortCode$ | async; as short) {
        <div class="max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
          <h1 class="text-3xl font-bold">Dashboard</h1>
          <h3 class="py-2 text-xl">
            You are watching details of Short <b>{{ short.Short }}</b
            >. It got clicked
            <b>{{ short.Count ?? 0 | numberToWords }}</b>
            {{ (short.Count ?? 0) === 1 ? 'time' : 'times' }} in total.
          </h3>
        </div>
        }
      </header>
      <div
        class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8 flex flex-col gap-4 w-full backdrop-blur-md rounded-3xl border border-gray-100 shadow-sm dark:border-gray-900"
      >
        <app-line-chart
          class="max-h-[30rem] min-h-[20rem]"
          chartTitle="Click histogram"
          [data$]="chartClickHistory$"
        ></app-line-chart>
      </div>
      } } } }
    </main>
  `,
})
export class ShortDetailsComponent {
  private shortId$: Observable<string>;
  shortDetails$: Observable<LoadState<ShortURL>>;
  shortCode$: Observable<ShortURL>;
  clickHistory$: Observable<LoadState<ShortClickCount[]>>;
  chartClickHistory$: Observable<DataPoint[]>;

  // Store navigation state
  private navigationState: ShortURL | null = null;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private dashboardService: DashboardService
  ) {
    // Capture navigation state in constructor
    const navigation = this.router.getCurrentNavigation();
    this.navigationState = navigation?.extras?.state?.['short'] || null;

    this.shortId$ = this.route.paramMap.pipe(
      map((params) => params.get('short') ?? ''),
      shareReplay(1)
    );

    this.shortDetails$ = this.shortId$.pipe(
      switchMap((id) => this.getShortDetails(id)),
      shareReplay(1)
    );

    this.shortCode$ = this.shortDetails$.pipe(
      filter(
        (s): s is { status: 'success'; data: ShortURL } =>
          s.status === 'success'
      ),
      map((s) => s.data ?? {})
    );

    this.clickHistory$ = this.shortId$.pipe(
      switchMap((id) => this.getClickHistory(id)),
      shareReplay(1)
    );

    this.chartClickHistory$ = this.clickHistory$.pipe(
      filter(
        (h): h is { status: 'success'; data: ShortClickCount[] } =>
          h.status === 'success'
      ),
      map((h) => this.transformToDataPoints(h.data))
    );
  }

  private getShortDetails(id: string): Observable<LoadState<ShortURL>> {
    // Use captured navigation state if available and matches ID
    if (this.navigationState && this.navigationState.Short === id) {
      return of({
        status: 'success' as const,
        data: this.navigationState,
      });
    }

    // Otherwise fetch from service
    return this.dashboardService.getShortDetails(id).pipe(
      map((data) => ({ status: 'success' as const, data })),
      catchError((error) =>
        of({
          status: 'error' as const,
          message: error.message || 'Failed to load short details',
        })
      ),
      startWith({ status: 'loading' as const })
    );
  }

  private getClickHistory(
    id: string
  ): Observable<LoadState<ShortClickCount[]>> {
    return this.dashboardService.clickHistoryByShort(id).pipe(
      map((data) => ({ status: 'success' as const, data: data })),
      catchError((error) =>
        of({
          status: 'error' as const,
          message: this.getErrorMessage(error),
        })
      ),
      startWith({ status: 'loading' as const })
    );
  }

  private transformToDataPoints(data: ShortClickCount[]): DataPoint[] {
    return data.map((item) => ({
      x: new Date(item.Timestamp ?? 0).getTime(),
      y: item.Count ?? 0,
    }));
  }

  private getErrorMessage(error: any): string {
    return error.status === 403
      ? 'Short does not exist, or you are not authorised to view it.'
      : error.message ?? 'An error occurred.';
  }
}
