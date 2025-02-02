import { Routes } from '@angular/router';
import { CollectionsListComponent } from './components/collections-list/collections-list.component';
import { SearchComponent } from './components/search/search.component';

export const routes: Routes = [
  { path: 'search', component: CollectionsListComponent },
  { path: 'search/:collection', component: SearchComponent },
  { path: '', redirectTo: '/search', pathMatch: 'full' },
];
