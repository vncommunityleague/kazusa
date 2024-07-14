import { Context, Namespace, SubjectSet } from "@ory/keto-namespace-types"

class User implements Namespace {}

class Role implements Namespace {
  related: {
    members: (User | Role)[]
  }
}

class Global implements Namespace {
  related: {
    editors: Role[]
  }

  permits = {
    tournament_edit: (ctx: Context): boolean => this.related.editors.includes(ctx.subject),
  }
}

class OsuTournament implements Namespace {
  related: {
    viewers: (User | SubjectSet<Role, "members">)[]
    editors: (User | SubjectSet<Role, "members">)[]

    mappool_viewers: (User | SubjectSet<Role, "members">)[]
    mappool_editors: (User | SubjectSet<Role, "members">)[]

    match_viewers: (User | SubjectSet<Role, "members">)[]
    match_editors: (User | SubjectSet<Role, "members">)[]

    global: Global[]
  }

  permits = {
    edit: (ctx: Context): boolean => 
      this.related.editors.includes(ctx.subject) ||
      this.related.global.traverse((p) => p.permits.tournament_edit(ctx)),

    mappool_edit: (ctx: Context): boolean => 
      this.permits.edit(ctx) ||
      this.related.mappool_editors.includes(ctx.subject),
  
    match_edit: (ctx: Context): boolean => 
      this.permits.edit(ctx) ||
      this.related.match_editors.includes(ctx.subject)
  }
}

class OsuMappool implements Namespace {
  related: {
    tournament: OsuTournament[]
  }

  permits = {
    edit: (ctx: Context): boolean =>
      this.related.tournament.traverse((p) => p.permits.mappool_edit(ctx))
  }
}

class OsuMatch implements Namespace {
  related: {
    tournament: OsuTournament[]
  }

  permits = {
    edit: (ctx: Context): boolean =>
      this.related.tournament.traverse((p) => p.permits.match_edit(ctx))
  }
}
