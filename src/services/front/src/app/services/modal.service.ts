// modal.service.ts
import {
  Injectable,
  ComponentRef,
  Type,
  EnvironmentInjector,
} from '@angular/core';
import { ModalHostComponent } from './../components/modal/modal-host/modal-host.component';

@Injectable({ providedIn: 'root' })
export class ModalService {
  private modalHost?: ModalHostComponent;
  private currentModal?: ComponentRef<any>;

  constructor(private injector: EnvironmentInjector) {}

  registerHost(host: ModalHostComponent): void {
    this.modalHost = host;
  }

  open<T>(component: Type<T>, inputs?: Record<string, any>): ComponentRef<T> {
    this.close();

    if (!this.modalHost) {
      throw new Error('Modal host not registered');
    }

    // ① Create in the host's ViewContainerRef
    const componentRef = this.modalHost.container.createComponent(component, {
      environmentInjector: this.injector,
    });

    // ② Pass inputs
    if (inputs) {
      Object.entries(inputs).forEach(([key, value]) => {
        componentRef.setInput(key, value);
      });
    }

    // ③ Show backdrop & keep ref
    this.modalHost.showBackdrop = true;
    this.currentModal = componentRef;
    return componentRef;
  }

  close(): void {
    if (this.modalHost) {
      // Remove any dynamic components + their DOM
      this.modalHost.container.clear();

      // Hide backdrop
      this.modalHost.showBackdrop = false;
    }
    this.currentModal = undefined;
  }
}
