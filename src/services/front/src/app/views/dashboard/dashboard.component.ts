import { Component } from '@angular/core';
import { DasherService, ShortURL } from '../../api';
import { Observable } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule],
  template: `
    <p>dashboard works!</p>
    @for (item of $allShorts | async; track item.Short) {
    <div>
      {{ item }}
    </div>
    }
  `,
  styles: ``,
})
export class DashboardComponent {
  public $allShorts: Observable<ShortURL[]>;
  constructor(private dasherService: DasherService) {
    this.$allShorts = this.dasherService.dasherShortAllGet();
  }
}
