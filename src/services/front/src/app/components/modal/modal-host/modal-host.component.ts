// modal-host.component.ts
import { CommonModule } from '@angular/common';
import { ModalService } from './../../../services/modal.service';
import {
  Component,
  ViewChild,
  ViewContainerRef,
  ChangeDetectorRef,
} from '@angular/core';

@Component({
  selector: 'app-modal-host',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div
      *ngIf="showBackdrop"
      class="fixed inset-0 backdrop-blur-xs z-[1000]"
      (click)="onBackdropClick()"
    ></div>

    <!-- Angular will stamp dynamic components here -->
    <ng-template #container></ng-template>
  `,
})
export class ModalHostComponent {
  showBackdrop = false;

  @ViewChild('container', { read: ViewContainerRef, static: true })
  container!: ViewContainerRef;

  constructor(
    private modalService: ModalService,
    private cdr: ChangeDetectorRef
  ) {
    this.modalService.registerHost(this);
  }

  onBackdropClick(): void {
    this.modalService.close();
  }
}
