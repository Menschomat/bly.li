import { Component, OnInit } from '@angular/core';
import { ButtonPrimaryComponent } from '../generic/button-primary/button-primary.component';
import { ShortnReq, ShortnService } from '../../core/api/v1';
import { FormsModule } from '@angular/forms';
import { BehaviorSubject } from 'rxjs';
import { CommonModule } from '@angular/common';
import { ConfigService } from '../../services/config.service';
import { URLService } from '../../services/url.service';

@Component({
  selector: 'app-url-input',
  standalone: true,
  imports: [ButtonPrimaryComponent, CommonModule, FormsModule],
  templateUrl: './url-input.component.html',
  styleUrl: './url-input.component.scss',
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
