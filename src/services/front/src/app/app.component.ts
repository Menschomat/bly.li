import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { UrlInputComponent } from './components/url-input/url-input.component';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';
import { UrlOutputComponent } from './components/url-output/url-output.component';

@Component({
    selector: 'app-root',
    imports: [
        RouterOutlet,
        UrlInputComponent,
        NavBarComponent,
        UrlOutputComponent,
    ],
    template: `
    <div class="flex flex-col h-full text-gray-800 dark:text-gray-200">
      <app-nav-bar></app-nav-bar>
      <div class="flex flex-1 flex-col items-center justify-center">
        <h2
          class="bg-animate text-center mb-5 text-6xl font-montserrat font-black leading-snug text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-pink-600 to-purple-600"
        >
          Shrink the Link, Elevate the Click
        </h2>
        <div class="flex flex-col gap-4">
          <app-url-input></app-url-input>
          <app-url-output></app-url-output>
        </div>
      </div>

      <div class="m-2 text-gray-700 dark:text-white flex justify-center">
        <div>©Mensch0 - 2024</div>
      </div>
    </div>
  `
})
export class AppComponent {
  title = 'front';
}
