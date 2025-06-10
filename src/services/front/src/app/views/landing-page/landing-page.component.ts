import { Component } from '@angular/core';
import { UrlInputComponent } from '../../components/url-input/url-input.component';
import { UrlOutputComponent } from '../../components/url-output/url-output.component';

@Component({
  selector: 'app-landing-page',
  standalone: true,
  imports: [UrlInputComponent, UrlOutputComponent],
  host: { class: 'flex flex-1 flex-col items-center justify-center' },
  template: `
    <h2
      class="bg-animate text-center mb-5 text-6xl font-montserrat font-black leading-snug text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-pink-600 to-purple-600"
    >
      &nbsp;&nbsp;{{ typewriterText }}<span class="animate-blink">|</span>
    </h2>
    <div class="flex flex-col gap-4">
      <app-url-input></app-url-input>
      <app-url-output></app-url-output>
    </div>
  `,
  styles: [`
    .animate-blink {
      animation: blink 1s step-end infinite;
    }
    @keyframes blink {
      0%, 100% { opacity: 1 }
      50% { opacity: 0 }
    }
  `],
})
export class LandingPageComponent {
  fullText = 'Shrink the Link, Elevate the Click';
  typewriterText = '';
  private i = 0;

  ngOnInit() {
    this.typeText();
  }

  typeText() {
    if (this.i < this.fullText.length) {
      this.typewriterText += this.fullText.charAt(this.i);
      this.i++;
      setTimeout(() => this.typeText(), 50); // speed here
    }
  }
}

