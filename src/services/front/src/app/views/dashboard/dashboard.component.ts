import { ShortnService } from '../../core/api/v1/api/shortn.service';
import { Component } from '@angular/core';
import { CardComponent } from '../../components/card/card.component';
import { UrlInputComponent } from '../../components/url-input/url-input.component';
import { UrlOutputComponent } from '../../components/url-output/url-output.component';
import { URLService } from '../../services/url.service';
import { Observable } from 'rxjs';
import { ShortURL } from '../../core/api/v1';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, CardComponent, UrlInputComponent, UrlOutputComponent],
  template: `
    <div>
      <app-url-input></app-url-input>
      <app-url-output></app-url-output>
    </div>
    <app-card>
      <div header>Your shorts...</div>
      <div body>
        <div class="relative overflow-x-auto shadow-md sm:rounded-lg">
          <table
            class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400"
          >
            <thead
              class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400"
            >
              <tr>
                <th scope="col" class="px-6 py-3">Short</th>
                <th scope="col" class="px-6 py-3">URL</th>
              </tr>
            </thead>
            <tbody>
              <tr
                class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                *ngFor="let item of ownUrls$ | async"
              >
                <td class="px-6 py-4">{{ item.Short }}</td>
                <td class="px-6 py-4">{{ item.URL }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </app-card>
  `,
  styles: `:host{
    @apply px-6 py-6 grid gap-4 grid-flow-col auto-cols-max content-center justify-items-center 
  }`,
})
export class DashboardComponent {
  constructor(
    private shortnService: ShortnService,
    private urlService: URLService
  ) {
    this.shortnService
      .allGet()
      .subscribe((a) => this.urlService.triggerNextOwnUrls(a));
  }
  get ownUrls$(): Observable<Array<ShortURL> | undefined> {
    return this.urlService.ownUrls$;
  }
}
