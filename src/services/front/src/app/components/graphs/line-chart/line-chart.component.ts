import { CommonModule } from '@angular/common';
import {
  Component,
  ViewChild,
  Input,
  OnInit,
  ChangeDetectionStrategy,
  effect,
} from '@angular/core';
import { resampleData } from '../../../utils/chart.utils';
import {
  ChartComponent,
  ApexAxisChartSeries,
  ApexChart,
  ApexXAxis,
  ApexDataLabels,
  ApexTitleSubtitle,
  ApexStroke,
  ApexGrid,
  ApexMarkers,
  NgApexchartsModule,
  ApexTheme,
} from 'ng-apexcharts';
import { Observable } from 'rxjs';
import { ThemeService } from '../../../services/theme.service';

export type ChartOptions = {
  series: ApexAxisChartSeries;
  chart: ApexChart;
  xaxis: ApexXAxis;
  dataLabels: ApexDataLabels;
  grid: ApexGrid;
  stroke: ApexStroke;
  title: ApexTitleSubtitle;
  markers: ApexMarkers;
  theme: ApexTheme;
  tooltip: ApexTooltip;
};

export interface DataPoint {
  x: number;
  y: number;
}

const COLOR_PALET = 'palette6';
@Component({
  selector: 'app-line-chart',
  imports: [CommonModule, NgApexchartsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: { class: 'flex-1 flex bg-white dark:bg-black' },
  template: `
    <apx-chart
      class=" flex-1"
      id="chart"
      #chart
      [series]="chartOptions.series"
      [chart]="chartOptions.chart"
      [xaxis]="chartOptions.xaxis"
      [dataLabels]="chartOptions.dataLabels"
      [grid]="chartOptions.grid"
      [stroke]="chartOptions.stroke"
      [title]="chartOptions.title"
      [markers]="chartOptions.markers"
      [theme]="chartOptions.theme"
      [tooltip]="chartOptions.tooltip"
    ></apx-chart>
  `,
  providers: [],
  styles: ``,
})
export class LineChartComponent implements OnInit {
  public chartOptions: ChartOptions;

  private _chartTitle: string = '';

  @ViewChild('chart') chart!: ChartComponent;

  @Input()
  public data$!: Observable<DataPoint[]>;

  @Input()
  set chartTitle(newValue: string) {
    this._chartTitle = newValue;
    //  this.chart.updateOptions(this.chartOptions);
  }

  get chartTitle(): string {
    return this._chartTitle ?? '';
  }

  constructor(private readonly themeService: ThemeService) {
    this.chartOptions = {
      series: [],
      chart: {
        height: '100%',
        width: '100%',
        type: 'line',
        zoom: {
          enabled: false,
        },
        background: 'transparent',
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        curve: 'smooth',
      },
      title: {
        text: this._chartTitle,
        align: 'left',
      },
      grid: {},
      xaxis: {
        type: 'datetime',
        labels: {
          datetimeFormatter: {
            year: 'yyyy',
            month: "MMM 'yy",
            day: 'dd MMM',
            hour: 'HH:mm',
          },
        },
      },
      tooltip: {
        x: {
          format: 'HH:mm dd.MM.yyyy',
        },
      },
      markers: {},
      theme: {
        palette: COLOR_PALET,
      },
    };

    effect(() => {
      const theme = this.themeService.theme() === 'lite' ? 'light' : 'dark';
      this.chartOptions.theme = {
        mode: theme,
        palette: COLOR_PALET,
      };
      this.chart?.updateOptions(this.chartOptions);
    });
  }

  ngOnInit(): void {
    // Resample at 1-hour intervals (3600000 ms)
    const tenMinMs = 10 * 60 * 1000;
    this.data$.subscribe((data) => {
      if (!data) return;
      this.chartOptions.series = [
        {
          name: 'Clicks',
          data: resampleData(data ?? [], tenMinMs).map((pt) => {
            return [pt.x, pt.y];
          }),
        },
      ];
      this.chart?.updateOptions(this.chartOptions);
    });
  }
}
