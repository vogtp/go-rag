import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { Document } from '../../services/collection-search.service';

@Component({
  selector: 'app-search-result-item',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './search-result-item.component.html',
  styleUrl: './search-result-item.component.css'
})
export class SearchResultItemComponent {

  @Input() doc: Document | undefined

}
