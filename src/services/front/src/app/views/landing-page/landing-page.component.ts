import { Component } from '@angular/core';
import { UrlInputComponent } from '../../components/url-input/url-input.component';
import { UrlOutputComponent } from '../../components/url-output/url-output.component';

@Component({
  selector: 'app-landing-page',
  imports: [UrlInputComponent, UrlOutputComponent],
  host: { class: 'flex flex-1 flex-col items-center justify-center' },
  template: `
    <h2
      class="bg-animate text-center mb-5 text-6xl font-montserrat font-black leading-snug text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-pink-600 to-purple-600"
    >
      Shrink the Link, Elevate the Click
    </h2>
    <div class="flex flex-col gap-4">
      <app-url-input></app-url-input>
      <app-url-output></app-url-output>
    </div>
  `,
  styles: ``,
})
export class LandingPageComponent {}
