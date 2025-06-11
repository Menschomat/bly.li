import { Component } from '@angular/core';
import { ButtonPrimaryComponent } from '../generic/button-primary/button-primary.component';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { URLService } from '../../services/url.service';
import { ShortnReq, ShortnService } from '../../api';

@Component({
  selector: 'app-url-input',
  standalone: true,
  imports: [ButtonPrimaryComponent, CommonModule, FormsModule],
  template: `
    <div
      class="mt-2 max-w-2xl px-6 py-4 rounded-3xl shadow-md 
                border border-gray-100 backdrop-blur-md 
                dark:border-gray-900"
    >
      <div class="flex flex-col">
        <a
          href="#"
          role="link"
          tabindex="0"
          class="px-2 text-lg text-gray-700 
                  dark:text-white"
        >
          Enter a long URL to create a short
        </a>

        <div class="flex flex-row items-center gap-2 mt-2">
          <input
            type="url"
            [(ngModel)]="urlInput"
            placeholder="https://example.url..."
            class="block w-full px-5 py-2.5 rounded-full 
                   border border-gray-200 bg-white 
                   text-gray-700 placeholder-gray-400/70 
                   focus:border-blue-400 focus:outline-none 
                   focus:ring focus:ring-blue-300 focus:ring-opacity-40 
                   dark:border-gray-600 dark:bg-transparent 
                   dark:text-gray-300 dark:placeholder-gray-500 
                   dark:focus:border-blue-300"
          />

          <app-button-primary (click)="requestShort()" class="flex-1 pl-1">
            {{ isLoading ? 'Shrinking...' : 'Shrink' }}
          </app-button-primary>
        </div>

        <p *ngIf="errorMessage" class="px-2 mt-2 text-sm text-red-500">
          {{ errorMessage }}
        </p>
      </div>
    </div>
  `,
})
export class UrlInputComponent {
  urlInput = '';
  isLoading = false;
  errorMessage: string | null = null;

  constructor(
    private readonly api: ShortnService,
    private readonly urlService: URLService
  ) {}

  requestShort() {
    if (!this.urlInput.trim()) return;

    this.isLoading = true;
    this.errorMessage = null;

    const request: ShortnReq = { Url: this.urlInput };

    this.api.shortnStorePost(request).subscribe({
      next: (response) => {
        this.urlService.triggerNextShort(response.Short);
        this.isLoading = false;
      },
      error: (err) => {
        this.errorMessage =
          err.message ?? 'Failed to shorten URL. Please try again.';
        this.isLoading = false;
      },
    });
  }
}
