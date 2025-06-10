import { Component, EventEmitter, Output, Input, OnInit } from '@angular/core';
import { NgClass } from '@angular/common';

// Define a type for our color configurations for better type safety
type ButtonColorConfig = {
  buttonGradient: string; // Classes for the outer button's gradient background
  spanGradient: string;   // Classes for the inner span's text gradient (MUST include direction like bg-gradient-to-r)
  hoverTextColor?: string; // Optional: if hover text color needs to change from default 'hover:text-white'
};

@Component({
  selector: 'app-button-primary',
  standalone: true,
  imports: [NgClass],
  template: `
    <button
      
      class="bg-animate p-0.5 rounded-full min-w-22"
      [ngClass]="getSelectedColorConfig().buttonGradient"
    >
      <div
        class="px-4 py-2 bg-white dark:bg-gray-900 rounded-full text-transparent hover:bg-transparent"
        [ngClass]="getSelectedColorConfig().hoverTextColor || 'hover:text-white'"
      >
        <span
          class="bg-animate block font-montserrat font-black leading-snug bg-clip-text"
          [ngClass]="getSelectedColorConfig().spanGradient"
        ><ng-content></ng-content></span>
      </div>
    </button>
  `,
  styles: []
})
export class ButtonPrimaryComponent implements OnInit {
  @Input() color: 'default' | 'green' | 'purple-pink' | 'blue-emerald' = 'default';
  @Output() clickEvent: EventEmitter<MouseEvent> = new EventEmitter<MouseEvent>();

  private colorConfigurations: Record<string, ButtonColorConfig> = {
    'default': {
      buttonGradient: 'from-indigo-400 via-pink-500 to-purple-500 bg-gradient-to-r',
      // CORRECTED: Added bg-gradient-to-r
      spanGradient: 'from-indigo-400 via-pink-500 to-purple-500 bg-gradient-to-r',
    },
    'green': { // This was your example, now correctly applied to spanGradient too
      buttonGradient: 'from-blue-400 to-emerald-400 bg-gradient-to-r',
      // CORRECTED: Added bg-gradient-to-r
      spanGradient: 'from-blue-400 to-emerald-400 bg-gradient-to-r',
    },
    'purple-pink': {
      buttonGradient: 'from-purple-500 to-pink-500 bg-gradient-to-r',
      // CORRECTED: Added bg-gradient-to-r
      spanGradient: 'from-purple-500 to-pink-500 bg-gradient-to-r',
    },
    'blue-emerald': { // Keeping the more descriptive name I added
      buttonGradient: 'from-blue-400 to-emerald-400 bg-gradient-to-r',
      // CORRECTED: Added bg-gradient-to-r
      spanGradient: 'from-blue-400 to-emerald-400 bg-gradient-to-r',
    }
    // Add more color configurations as needed, ensuring spanGradient includes a direction class
  };

  ngOnInit(): void {
    if (!this.colorConfigurations[this.color]) {
      console.warn(`ButtonPrimaryComponent: Unknown color '${this.color}'. Defaulting to 'default'.`);
      this.color = 'default';
    }
  }



  public getSelectedColorConfig(): ButtonColorConfig {
    return this.colorConfigurations[this.color] || this.colorConfigurations['default'];
  }
}