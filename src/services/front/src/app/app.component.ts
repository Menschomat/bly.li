import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { UrlInputComponent } from './components/url-input/url-input.component';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';
import { UrlOutputComponent } from './components/url-output/url-output.component';
import { MainComponent } from './views/main/main.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    UrlInputComponent,
    NavBarComponent,
    UrlOutputComponent,
    MainComponent,
  ],
  template: `
    <div class="flex flex-col h-full text-gray-800 dark:text-gray-200">
      <app-nav-bar></app-nav-bar>
      <div class="flex flex-1 flex-col items-center justify-center">
        <router-outlet></router-outlet>
      </div>

      <div class="m-2 text-gray-700 dark:text-white flex justify-center">
        <div>©Mensch0 - 2024</div>
      </div>
    </div>
  `,
})
export class AppComponent {
  title = 'front';
}
