import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { NgxEchartsDirective, provideEchartsCore } from 'ngx-echarts';
// import echarts core
import * as echarts from 'echarts/core';
// import necessary echarts components
import { BarChart, LineChart } from 'echarts/charts';
import { GridComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import { Observable, of } from 'rxjs';
import { EChartsOption } from 'echarts';
echarts.use([BarChart, LineChart, GridComponent, CanvasRenderer]);

@Component({
  selector: 'app-line-chart',
  imports: [CommonModule, NgxEchartsDirective],
  template: `
    <div
      echarts
      [options]="options$ | async"
      [loading]="loading"
      class="demo-chart"
    ></div>
  `,
  providers: [provideEchartsCore({ echarts })],
  styles: ``,
})
export class LineChartComponent {
  loading = true;
  @Input()
  options$: Observable<EChartsOption> = of({});
}
