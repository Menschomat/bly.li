import { Routes } from '@angular/router';
import { LandingPageComponent } from './views/landing-page/landing-page.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';

export const routes: Routes = [
  {
    // Standardroute: Umleitung auf '/home'
    path: '',
    redirectTo: 'front',
     pathMatch: 'full'
  },
  {
    path: 'front',
    component: LandingPageComponent,
  },
  {
    path: 'front/dash',
    component: DashboardComponent,
  },
];
