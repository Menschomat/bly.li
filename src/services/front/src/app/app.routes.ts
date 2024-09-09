import { Routes } from '@angular/router';

import { MainComponent } from './views/main/main.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';

export const routes: Routes = [
  {
    path: '',
    component: MainComponent,
  },
  {
    path: 'front',
    component: MainComponent,
  },
  {
    path: 'dash',
    component: DashboardComponent,
  },
];
