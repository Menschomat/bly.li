import { Routes } from '@angular/router';
import { LandingPageComponent } from './views/landing-page/landing-page.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';

export const routes: Routes = [
  {
    path: 'front',
    component: LandingPageComponent,
  },
  {
    path: 'dash',
    component: DashboardComponent,
  },
  {
    // Standardroute: Umleitung auf '/home'
    path: '',
    redirectTo: 'front',
    pathMatch: 'full',
  },
  {
    // Standardroute: Umleitung auf '/home'
    path: '**',
    redirectTo: 'front',
    pathMatch: 'full',
  },
];
