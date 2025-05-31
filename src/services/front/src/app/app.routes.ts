import { Routes } from '@angular/router';
import { LandingPageComponent } from './views/landing-page/landing-page.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';

export const routes: Routes = [
  {
    path: 'main',
    component: LandingPageComponent,
  },
  {
    path: 'dash',
    component: DashboardComponent,
  },
  {
    // Standardroute: Umleitung auf '/home'
    path: '',
    redirectTo: 'main',
    pathMatch: 'full',
  },
  {
    // Standardroute: Umleitung auf '/home'
    path: '**',
    redirectTo: 'main',
    pathMatch: 'full',
  },
];
