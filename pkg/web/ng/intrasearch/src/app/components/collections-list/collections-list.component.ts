import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { CollectionListService, CollectionRequest } from '../../services/collection-list.service';

import { MatListModule } from '@angular/material/list';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-collections-list',
  standalone: true,
  imports: [CommonModule, MatListModule, RouterLink],
  templateUrl: './collections-list.component.html',
  styleUrl: './collections-list.component.css'
})
export class CollectionsListComponent {
  collectionRequest: CollectionRequest | undefined;

  constructor(private collectionService: CollectionListService) {}

  
  loadCollections() {
    this.collectionService.getCollections()
      .subscribe(data => {
        this.collectionRequest = data
      }
      );
  }
  
 ngOnInit() {
   this.loadCollections()
 }
}
