import { Component, HostListener } from '@angular/core';
import { ModalService } from '../../../services/modal.service';

@Component({ template: '' })
export abstract class BaseModalComponent {
  constructor(protected modalService: ModalService) {}

  @HostListener('document:keydown.escape')
  handleEscape(): void {
    this.close();
  }

  close(): void {
    this.modalService.close();
  }
}
