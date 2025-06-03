import { Component } from '@angular/core';
import { map, Observable } from 'rxjs';
import { URLService } from '../../services/url.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ConfigService } from '../../services/config.service';
@Component({
  selector: 'app-url-output',
  imports: [CommonModule, FormsModule],
  template: `
    <div
      *ngIf="outUrl$ | async as outUrl"
      class=" max-w-2xl px-8 py-4  backdrop-blur-md bg-white/30 shadow dark:bg-gray-800/50 rounded-lg"
    >
      <div class="flex flex-row items-center gap-2">
        <input
          readonly
          disabled
          type="text"
          [value]="outUrl"
          placeholder="https://bly.li/...."
          class="block  w-full placeholder-gray-400/70 dark:placeholder-gray-500 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-gray-700 focus:border-blue-400 focus:outline-none focus:ring focus:ring-blue-300 focus:ring-opacity-40 dark:border-gray-600 dark:bg-gray-900 dark:text-gray-300 dark:focus:border-blue-300"
        />
        <button
          (click)="copyToClipboard(outUrl)"
          class="bg-animate p-0.5 rounded-lg bg-gradient-to-r from-blue-400 to-emerald-400"
        >
          <div
            class="px-4 pt-1.5 pb-2.5 bg-white dark:bg-gray-900 rounded-md text-transparent hover:bg-transparent hover:text-white"
          >
            <span
              class="bg-animate block font-montserrat font-black leading-snug bg-clip-text bg-gradient-to-r from-blue-400 via-emerald-400 to-cyan-400"
            >
              <!--<i class="text-2xl fa-regular fa-clipboard"></i>-->
              <p>copy</p>
            </span>
          </div>
        </button>
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
