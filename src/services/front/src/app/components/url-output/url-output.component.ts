import { Component } from '@angular/core';
import { map, Observable } from 'rxjs';
import { URLService } from '../../services/url.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ConfigService } from '../../services/config.service';
import { ButtonPrimaryComponent } from '../generic/button-primary/button-primary.component';
@Component({
  selector: 'app-url-output',
  imports: [ButtonPrimaryComponent, CommonModule, FormsModule],
  template: `
    <div
      *ngIf="outUrl$ | async as outUrl"
      class="mt-0 max-w-2xl px-6 py-4 rounded-3xl shadow-md 
                border border-gray-100 backdrop-blur-md 
                dark:border-gray-900"
    >
      <div class="flex flex-row items-center gap-2">
        <input
          readonly
          disabled
          type="text"
          [value]="outUrl"
          placeholder="https://bly.li/...."
          class="block w-full px-5 py-2.5 rounded-full 
                   border border-gray-200 bg-white 
                   text-gray-700 placeholder-gray-400/70 
                   focus:border-blue-400 focus:outline-none 
                   focus:ring focus:ring-blue-300 focus:ring-opacity-40 
                   dark:border-gray-600 dark:bg-transparent 
                   dark:text-gray-300 dark:placeholder-gray-500 
                   dark:focus:border-blue-300"
        />
        <app-button-primary
          [color]="'green'"
          (click)="copyToClipboard(outUrl)"
          class="flex-1 pl-1"
          >copy</app-button-primary
        >
      </div>
    </div>
  `,
  styles: ``,
})
export class UrlOutputComponent {
  isCopied: boolean = false;
  constructor(
    private readonly urlService: URLService,
    private readonly config: ConfigService
  ) {}
  get outUrl$(): Observable<string | undefined> {
    return this.urlService.lastShortUrl$.pipe(
      map((short) =>
        short
          ? `${window.location.origin}${
              this.config.getConfig().blowUpPath
            }/${short}`
          : undefined
      )
    );
  }

  async copyToClipboard(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      this.isCopied = true;
    } catch (err) {
      console.error('Failed to copy: ', err);
    }

    // Remove the "Copied!" message after a short delay
    setTimeout(() => {
      this.isCopied = false;
    }, 2000);
  }

  public dismiss() {
    this.urlService.clearLastShort();
  }
}
