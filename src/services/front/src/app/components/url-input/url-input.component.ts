import { Component, OnInit } from '@angular/core';
import { ButtonPrimaryComponent } from '../generic/button-primary/button-primary.component';
import { ShortnReq, ShortnService } from '../../core/api/v1';
import { FormsModule } from '@angular/forms';
import { BehaviorSubject } from 'rxjs';
import { CommonModule } from '@angular/common';
import { URLService } from '../../services/url.service';

@Component({
  selector: 'app-url-input',
  standalone: true,
  imports: [ButtonPrimaryComponent, CommonModule, FormsModule],
  template: `
    <div
      class="mt-2 max-w-2xl px-8 py-4 backdrop-blur-md bg-white/30 shadow dark:bg-gray-800/50 rounded-lg"
    >
      <div class=" flex flex-col">
        <a
          href="#"
          class="text-xl font-bold text-gray-700 dark:text-white"
          tabindex="0"
          role="link"
          >Enter a long URL to create your short:</a
        >
        <div class="flex flex-row items-center">
          <input
            type="text"
            [(ngModel)]="shortInputValue"
            placeholder="https://example.url..."
            class="block mt-2 w-full placeholder-gray-400/70 dark:placeholder-gray-500 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-gray-700 focus:border-blue-400 focus:outline-none focus:ring focus:ring-blue-300 focus:ring-opacity-40 dark:border-gray-600 dark:bg-gray-900 dark:text-gray-300 dark:focus:border-blue-300"
          />
          <app-button-primary
            (click)="requestShort()"
            class="flex-1 mt-2 pl-1"
          >shrink</app-button-primary>
        </div>
      </div>
      <a
        [href]="baseUrl + '/' + lastShort"
        target="_blank"
        *ngIf="lastShort$ | async as lastShort"
      >
        {{ baseUrl + '/' + lastShort }}
      </a>
    </div>
  `,
})
export class UrlInputComponent implements OnInit {
  public shortInputValue: string = '';
  public baseUrl: string = window.location.origin;
  private lastShortSubj: BehaviorSubject<string | undefined> =
    new BehaviorSubject<string | undefined>(undefined);
  constructor(private api: ShortnService, private urlService: URLService) {}
  ngOnInit(): void {}
  requestShort() {
    this.api
      .storePost({ Url: this.shortInputValue } as ShortnReq)
      .subscribe((a) => this.urlService.triggerNextShort(a.Short));
  }
  get lastShort$() {
    return this.lastShortSubj.asObservable();
  }
}
