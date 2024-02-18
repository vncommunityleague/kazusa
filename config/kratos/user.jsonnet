function(ctx) {
  id : ctx.identity.id,
  username: ctx.identity.traits.username,
  email : ctx.identity.traits.email,
  name: if "name" in ctx.identity.traits then ctx.identity.traits.name else null,
  default_picture: if ctx.identity.metadata_public != null && "default_picture" in ctx.identity.metadata_public then ctx.identity.metadata_public.default_picture else null,
  discord_id: if ctx.identity.metadata_public != null && "discord_id" in ctx.identity.metadata_public then ctx.identity.metadata_public.discord_id else null,
}
