import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIcon } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatToolbarModule } from '@angular/material/toolbar';
import { RouterOutlet } from '@angular/router';
import { SpaceListButtonComponent } from './components/space-list-button/space-list-button.component';
import { SettingsButtonComponent } from "./components/settings-button/settings-button.component";
@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    MatToolbarModule,
    MatIcon,
    MatMenuModule,
    MatButtonModule,
    SpaceListButtonComponent,
    SettingsButtonComponent
],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent {
  title = 'intrasearch';
}
