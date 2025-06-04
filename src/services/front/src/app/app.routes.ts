import { Routes } from '@angular/router';
import { LandingPageComponent } from './views/landing-page/landing-page.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';
import { ShortDetailsComponent } from './views/dashboard/short-details/short-details.component';

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
    path: 'dash/detail/:short',
    component: ShortDetailsComponent,
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
